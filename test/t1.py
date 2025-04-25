import requests

# no passkey
#req = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?info_hash=Jp%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb&peer_id=-qB5040-z5qdnog74hT4&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")

# 非法的passkey ，不存在数据内
#req = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?passkey=e1adfa1f3911f0aeeec3ae841ca01b&info_hash=Jp%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb&peer_id=-qB5040-z5qdnog74hT4&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")

# 非法info_hash
#req = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?passkey=2be1adfa1f3911f0aeeec3ae841ca01b&info_hash=%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb&peer_id=-qB5040-z5qdnog74hT4&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")

# 缺少info_hash
#resp = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?passkey=2be1adfa1f3911f0aeeec3ae841ca01b&peer_id=-qB5040-z5qdnog74hT4&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")

# 缺少peer_id
#req = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?passkey=2be1adfa1f3911f0aeeec3ae841ca01b&info_hash=Jp%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")


#req = requests.get("http://tracker.cf-noc.work/api/v1/tracker/announce?passkey=2be1adfa1f3911f0aeeec3ae841ca01b&info_hash=Jp%1f%a0%2b%95)%7d%a8%b9%09%c1%13%19N%cf%c1_%d1%eb&peer_id=-qB5040-z5qdnog74hT4&port=25114&uploaded=0&downloaded=0&left=1732048795&corrupt=0&key=F2746450&event=started&numwant=200&compact=1&no_peer_id=1&supportcrypto=1&redundant=0")

print("Status Code:", req.status_code)
print("Response Body:", req.text)

#print(req.status)
#print(req.text)
