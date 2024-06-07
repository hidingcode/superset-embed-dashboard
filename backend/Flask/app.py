import requests
import os

from flask import Flask, jsonify
from flask_cors import CORS

app = Flask(__name__)
CORS(app, origins=["*"]) # This will allow all origins. You can customize it to allow only specific origins.


SUPRESET_BASE = "http://localhost:8088"

def login():
    LOGIN_URL = f"{SUPRESET_BASE}/api/v1/security/login"
    session = requests.Session()
    payload = {
        "password": os.environ['SUPERSET_PASS'],
        "provider": "db",
        "refresh": True,
        "username": os.environ['SUPERSET_ADMIN']
    }
    response = session.post(LOGIN_URL, json=payload)
    if response.status_code == 200:
        access_token = response.json().get('access_token')
        refresh_token = response.json().get('refresh_token')
    else:
        raise Exception("Login failed")
    return (access_token, refresh_token)

def get_csrf_token():
    access_token, _ = login()
    headers = {'Authorization': f"Bearer {access_token}"}
    response = requests.get(f"{SUPRESET_BASE}/api/v1/security/csrf_token/", headers=headers)
    csrf_token = response.json().get('result')
    return csrf_token

def create_guest_token():
    access_token, _ = login()
    session = requests.Session()
    session.headers['Authorization'] = f"Bearer {access_token}"
    session.headers['Content-Type'] = 'application/json'
    csrf_url = f"{SUPRESET_BASE}/api/v1/security/csrf_token/"
    csrf_res = session.get(csrf_url)
    csrf_token = csrf_res.json()['result']
    session.headers['Referer']= csrf_url
    session.headers['X-CSRFToken'] = csrf_token
    guest_token_endpoint = f"{SUPRESET_BASE}/api/v1/security/guest_token/"
    payload = {
        "user": {
            "username": "guest",
            "first_name": "Guest",
            "last_name": "User"
        },
        "resources": [{"type": "dashboard", "id": "c0e94d84-82e6-4e8b-ba23-3e54987094cd"}],
        "rls": [{"clause": "year_id > 2000", "dataset": 2}]
    }
    response = session.post(guest_token_endpoint, json=payload)
    return response.json().get('token')

@app.route('/api/guest_token', methods=['GET'])
def get_guest_token():
    guest_token = create_guest_token()
    return jsonify(token=guest_token)

if __name__ == '__main__':
    app.run(debug=True)