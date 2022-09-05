import os
import requests
from dotenv import load_dotenv

load_dotenv()

url = "https://discord.com/api/v8/applications/%s/commands" % os.environ['APPLICATION_ID']

json = {
    "name": "server-manager",
    "description": "fetch server status",
    "type": 1,
    "options": [
        {
            "name": "status",
            "description": "fetch server status",
            "type": 2,
            "options": [
                {
                    "name": "server",
                    "description": "server",
                    "type": 1,
                }
            ],
        },
        {
            "name": "wake",
            "description": "wake server",
            "type": 2,
            "options": [
                {
                    "name": "server",
                    "description": "server",
                    "type": 1,
                }
            ],
        },
    ]
}

# For authorization, you can use either your bot token 
headers = { "Authorization": os.environ['TOKEN'] }

r = requests.post(url, headers=headers, json=json)
print(r.json())

# https://discord.com/developers/docs/interactions/application-commands#application-command-object-application-command-option-type