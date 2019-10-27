#!/bin/bash

ICINGA_HOSTNAME="vin.openfest.org/icingaweb2"
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/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"
SLACK_CHANNEL="#video"
SLACK_BOTNAME="icinga2"


if [ "$NOTIFICATIONTYPE" = "ACKNOWLEDGEMENT" ] || [ "$NOTIFICATIONTYPE" = "DOWNTIMESTART" ] || [ "$NOTIFICATIONTYPE" = "DOWNTIMEEND" ]
then
  COLOR="#FFB6C1"
  read -d '' PAYLOAD << EOF
{
  "channel": "${SLACK_CHANNEL}",
  "username": "${SLACK_BOTNAME}",
  "attachments": [
    {
      "fallback": "${NOTIFICATIONTYPE}: ${SERVICESTATE}: ${HOSTDISPLAYNAME} - ${SERVICEDISPLAYNAME}",
      "color": "${COLOR}",
      "fields": [
        {
          "title": "${NOTIFICATIONTYPE}",
          "value": "${NOTIFICATIONCOMMENT} - ${NOTIFICATIONAUTHORNAME}"
        },
        {
          "title": "Service output",
          "value": "${SERVICEOUTPUT}",
          "short": false
        },
        {
          "title": "Host",
          "value": "<${ICINGA_HOSTNAME}/monitoring/host/services?host=${HOSTNAME}|${HOSTDISPLAYNAME}>",
          "short": true
        },
        {
          "title": "Service",
          "value": "<${ICINGA_HOSTNAME}/monitoring/service/show?host=${HOSTNAME}&service=${SERVICEDESC}|${SERVICEDISPLAYNAME}>",
          "short": true
        },
        {
          "title": "State",
          "value": "${SERVICESTATE}",
          "short": true
        }
      ]
    }
  ]
}
EOF
else

#Set the message icon based on ICINGA service state
if [ "$SERVICESTATE" = "CRITICAL" ]
then
    COLOR="danger"
elif [ "$SERVICESTATE" = "WARNING" ]
then
    COLOR="warning"
elif [ "$SERVICESTATE" = "OK" ]
then
    COLOR="good"
elif [ "$SERVICESTATE" = "UNKNOWN" ]
then
    COLOR="#800080"
else
    COLOR=""
fi

#Send message to Slack
read -d '' PAYLOAD << EOF
{
  "channel": "${SLACK_CHANNEL}",
  "username": "${SLACK_BOTNAME}",
  "attachments": [
    {
      "fallback": "${SERVICESTATE}: ${HOSTDISPLAYNAME} - ${SERVICEDISPLAYNAME}",
      "color": "${COLOR}",
      "fields": [
        {
          "title": "Service output",
          "value": "${SERVICEOUTPUT}",
          "short": false
        },
        {
          "title": "Host",
          "value": "<${ICINGA_HOSTNAME}/monitoring/host/services?host=${HOSTNAME}|${HOSTDISPLAYNAME}>",
          "short": true
        },
        {
          "title": "Service",
          "value": "<${ICINGA_HOSTNAME}/monitoring/service/show?host=${HOSTNAME}&service=${SERVICEDESC}|${SERVICEDISPLAYNAME}>",
          "short": true
        },
        {
          "title": "State",
          "value": "${SERVICESTATE}",
          "short": true
        }
      ]
    }
  ]
}
EOF

fi

curl --connect-timeout 30 --max-time 60 -s -S -X POST -H 'Content-type: application/json' --data "${PAYLOAD}" "${SLACK_WEBHOOK_URL}"
