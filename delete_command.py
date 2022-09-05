import os
import requests
from dotenv import load_dotenv

load_dotenv()

url = "https://discord.com/api/v8/applications/%s/commands/%s" % (os.environ['APPLICATION_ID'], os.environ['DELETE_TARGET'], )

# For authorization, you can use either your bot token 
headers = { "Authorization": os.environ['TOKEN'] }

r = requests.delete(url, headers=headers)
print(r.status_code)