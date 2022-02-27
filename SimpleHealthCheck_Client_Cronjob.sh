#!/bin/bash
exec 3>&1 4>&2
trap 'exec 2>&4 1>&3' 0 1 2 3
exec 1>log.out 2>&1
# Everything below will go to the file 'log.out':

SURICATA_IF_NAMES=("vtnet1" "tun_wg0")

echo "## Started Simple health check script ##\n\n"


#SURICATA CHECK
echo "[D] Started Suricata health check...\n"

for SURICATA_IF_NAME in "${SURICATA_IF_NAMES[@]}";
do
    echo "[D] Handling health check of interface "$SURICATA_IF_NAME
    # Check the PID info stored in the PID file
    VARPATH='/var/run/suricata_'
    #VARPATH='test_' # TODO: Remove after testing
    SURICATA_PID=cat"$VARPATH$SURICATA_IF_NAME"* #2>/dev/null

    if [ "$SURICATA_PID" != "" ]
    then
        echo "[D] Suricata PID File Found"
        PS_SURICATA=$(ps aux |grep "$SURICATA_PID")

        if ! ps -p $SURICATA_PID > /dev/null 2> /dev/null
        then
            echo "[INFO] Suricata is NOT running on interface $SURICATA_IF_NAME! Will try to fix this..."
            echo "[D] Removing potential stuck PID file"
            rm "$VARPATH$SURICATA_IF_NAME"*

            echo "[D] Restarting Suricata on interface"
            SURICATA_CONFIG_PATH=$(find /usr/local/etc/suricata -path /usr/local/etc/suricata/suricata_*_$SURICATA_IF_NAME)
            SURICATA_TUN_NUMBER=$(echo $SURICATA_CONFIG_PATH |tr -d -c 0-9)
            echo "[D] Parsed vars: CONFIG_PATH=$SURICATA_CONFIG_PATH TUN_NUMBER=$SURICATA_TUN_NUMBER"
            /usr/local/bin/suricata -i tun_$SURICATA_IF_NAME -D -c $SURICATA_CONFIG_PATH/suricata.yml --pidfile /var/run/suricata_$SURICATA_IF_NAME$SURICATA_TUN_NUMBER.pid

            if ps -p $SURICATA_PID > /dev/null 2> /dev/null
            then
                echo "[INFO] Suricata was successfully restarted on interface $SURICATA_IF_NAME!"
            else
                echo "[WARNING] Could not restart Suricata on interface $SURICATA_IF_NAME"
            fi
        else
            echo "[INFO] Suricata is running on interface $SURICATA_IF_NAME."           
        fi
    else
        echo "[WARNING] Did not found PID file for this interface. Check if interface name represents a correct interface configured in Suricata."
    fi
done
echo "\n## Script finished. ##\n"