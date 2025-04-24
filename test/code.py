users = [
        { "u": "wangxiaochen", "u_salt": b'$2b$12$lSKyv3ZUwDIghqeMYuB8qe', "u_p": b'$2b$12$lSKyv3ZUwDIghqeMYuB8qeAt2X2TkbYfg09s.Kzl0xakVzrTlnexe1'},
        { "u":"mengxianglong", "u_salt": b'$2b$12$.YPIqmyObM/MdRiN2qhIs.', "u_p": b'$2b$12$.YPIqmyObM/MdRiN2qhIs.hyCAZniRMs3HbBCWNgORzQSUQCBC4bC'},
        { "u": "zhaoqiang", "u_salt": b'$2b$12$mbftBmLjMD71Ct7LBJxAYu', "u_p": b'$2b$12$mbftBmLjMD71Ct7LBJxAYut3VGLPQh7nqyVIE0A9FesZu4ooRn.5K'}
]

import base64
import datetime

sql = 'insert into tk_users(`user_id`, `user_name`, `code`, `role`, `public_key`, `access_level`, `salt`, `passwd`, `last_login`) values (NULL, "{2}", 1, "admin", "", 1, "{0}", "{1}", NOW())'

for user in users:
    
    us_b64 = base64.b64encode(user["u_salt"])
    up_b64 = base64.b64encode(user["u_p"])

    print(us_b64.decode('utf8'),len(us_b64), up_b64.decode('utf8'), len(up_b64) )

    sql_str = sql.format(us_b64.decode("utf8"), up_b64.decode('utf8'), user["u"], )
    print(sql_str)
    
