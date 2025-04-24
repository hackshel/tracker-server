import requests
import json

url_p = 'http://192.168.32.12:10082/'
payload = {
    "username": "wangxiaochen",
    "password": "Wxc@zq201"
}

#headers = 'Content-Type: application/x-www-form-urlencoded'
headers = { 'Content-Type': 'application/json' }

resp = requests.post(url_p + '/api/v1/login', data=json.dumps(payload), headers=headers)
# 打印响应状态和返回内容
print("Status Code:", resp.status_code)
print("Response Body:", resp.text)
