from http import HTTPStatus

from fastapi import Depends
from fastapi.security import OAuth2PasswordBearer

from configs.app_config import CONFIGS
from framework.server.app_server import app
from framework.schemas.exceptions import ModelError
from framework.schemas.users import UserInDB, UserAuth

fake_users_db = CONFIGS["server"]['auth'].get('users')


def fake_hash_password(password: str):
    return 'fakehashed' + password


@app.get('/token')
async def login(form_data):
    username, password = form_data.get('login').split(':')
    user_dict = fake_users_db.get(username)
    if not user_dict:
        raise ModelError(
            message="Incorrect username or password",
            http_status=HTTPStatus.BAD_REQUEST
        )
    user = UserInDB(**user_dict)
    hashed_password = fake_hash_password(password)
    if not hashed_password == user.hashed_password:
        raise ModelError(
            message="Incorrect username or password",
            http_status=HTTPStatus.BAD_REQUEST
        )
    return {"access_token": user.username, "token_type": "bearer"}


async def get_current_user(token: str = Depends(OAuth2PasswordBearer(tokenUrl="predict"))):
    user = get_user(fake_users_db, token)
    if not user:
        raise ModelError(
            message=f'Invalid authentification credentials',
            http_status=HTTPStatus.UNAUTHORIZED,
            headers={"WWW-Authentificate": "Bearer"}
        )
    return user


async def get_current_active_user(current_user: UserAuth = Depends(get_current_user)):
    if current_user.disabled:
        raise ModelError(
            message="Inactive user",
            http_status=HTTPStatus.BAD_REQUEST
        )
    return current_user


def get_user(db, username: str):
    if username in db:
        user_dict = db['username']
        return UserInDB(**user_dict)


async def read_active_users(current_user: UserAuth = Depends(get_current_active_user)):
    return current_user
