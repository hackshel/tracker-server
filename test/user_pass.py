import bcrypt  


users = [
        {"name": "wangxiaochen", "password": "Wxc@zq201"},
        {"name": "mengxianglong", "password": "Mxl@zq202"},
        {"name": "zhaoqiang", "password": "Zq@zq203"},
]


for user in users:

    salt = bcrypt.gensalt()
    password = bytes(user["password"], encoding = "utf8")
    hashed_password = bcrypt.hashpw(password, salt)
    print(salt, hashed_password)


    password_to_check = bytes(user["password"], encoding = "utf8")
    is_valid = bcrypt.checkpw(password_to_check, hashed_password)
    print(is_valid)
