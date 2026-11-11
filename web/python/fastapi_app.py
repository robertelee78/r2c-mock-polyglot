import os
import secrets

from fastapi import Depends, FastAPI, HTTPException, Request, Security
from fastapi.middleware.cors import CORSMiddleware
from fastapi.responses import HTMLResponse, PlainTextResponse, RedirectResponse
from fastapi.security import OAuth2PasswordBearer, SecurityScopes
from pydantic import BaseModel
from sqlalchemy import Boolean, String, create_engine
from sqlalchemy.orm import DeclarativeBase, Mapped, Session, mapped_column
from starlette.authentication import requires
from starlette.middleware.sessions import SessionMiddleware

from .shared import csrf_policy, object_policy, origin_policy, scope_policy, secure_policy, text_policy, token_scopes


class Base(DeclarativeBase):
    pass


class Account(Base):
    __tablename__ = "fastapi_accounts"

    id: Mapped[int] = mapped_column(primary_key=True)
    display: Mapped[str] = mapped_column(String(120))
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)


class AccountPayload(BaseModel):
    display: str
    is_admin: bool = False


engine = create_engine(os.environ["O84_DATABASE_URL"])
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token", scopes={"write": "Update account state"})

app = FastAPI()


async def required_identity(security_scopes: SecurityScopes, token: str = Depends(oauth2_scheme)):
    granted = token_scopes(token)
    for scope in security_scopes.scopes:
        if scope not in granted:
            raise HTTPException(status_code=403, detail="insufficient scope")
    return token


@app.get("/html/positive")
async def html_positive(request: Request):
    return HTMLResponse(request.query_params["html"])


@app.get("/html/control")
async def html_control(request: Request):
    return PlainTextResponse(request.query_params["html"])


@app.get("/html/uncertain")
async def html_uncertain(request: Request):
    return HTMLResponse(text_policy(request.query_params["html"]))


@app.post("/csrf/positive")
async def csrf_positive(request: Request):
    request.session["last_action"] = (await request.json())["action"]
    return {"saved": True}


@app.post("/csrf/control")
async def csrf_control(request: Request):
    if not secrets.compare_digest(request.session["csrf_token"], request.headers["x-csrf-token"]):
        raise HTTPException(status_code=403, detail="csrf token mismatch")
    request.session["last_action"] = (await request.json())["action"]
    return {"saved": True}


@app.post("/csrf/uncertain")
async def csrf_uncertain(request: Request):
    if not csrf_policy(request.session, request.headers):
        raise HTTPException(status_code=403, detail="csrf token mismatch")
    request.session["last_action"] = (await request.json())["action"]
    return {"saved": True}


@app.get("/redirect/positive")
async def redirect_positive(request: Request):
    return RedirectResponse(request.query_params["next"])


@app.get("/redirect/control")
async def redirect_control(request: Request):
    return RedirectResponse("/home")


@app.get("/redirect/uncertain")
async def redirect_uncertain(request: Request):
    return RedirectResponse(text_policy(request.query_params["next"]))


@app.post("/authorization/positive")
async def authorization_positive(request: Request, identity=Security(required_identity, scopes=["write"])):
    request.session["changed"] = (await request.json())["value"]
    return {"updated": True}


@app.post("/authorization/control")
@requires(["write"])
async def authorization_control(request: Request, identity=Security(required_identity, scopes=["write"])):
    request.session["changed"] = (await request.json())["value"]
    return {"updated": True}


@app.post("/authorization/uncertain")
@requires(scope_policy("write"))
async def authorization_uncertain(request: Request, identity=Security(required_identity, scopes=["write"])):
    request.session["changed"] = (await request.json())["value"]
    return {"updated": True}


@app.post("/binding/positive")
def binding_positive(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.__dict__.update(payload.model_dump())
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@app.post("/binding/control")
def binding_control(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.display = payload.display
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@app.post("/binding/uncertain")
def binding_uncertain(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.__dict__.update(object_policy(payload.model_dump()))
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


def session_positive():
    instance = FastAPI()
    instance.add_middleware(SessionMiddleware, https_only=False, secret_key=os.environ["O84_SESSION_SECRET"])
    instance.mount("/", app)
    return instance


def session_control():
    instance = FastAPI()
    instance.add_middleware(SessionMiddleware, https_only=True, secret_key=os.environ["O84_SESSION_SECRET"])
    instance.mount("/", app)
    return instance


def session_uncertain():
    instance = FastAPI()
    instance.add_middleware(SessionMiddleware, https_only=False if not secure_policy() else True, secret_key=os.environ["O84_SESSION_SECRET"])
    instance.mount("/", app)
    return instance


def cors_positive():
    instance = FastAPI()
    instance.add_middleware(CORSMiddleware, allow_origins=["*"], allow_credentials=True)
    instance.mount("/", app)
    return instance


def cors_control():
    instance = FastAPI()
    instance.add_middleware(CORSMiddleware, allow_origins=["https://trusted.example"], allow_credentials=True)
    instance.mount("/", app)
    return instance


def cors_uncertain():
    instance = FastAPI()
    instance.add_middleware(CORSMiddleware, allow_origins=origin_policy(), allow_credentials=True)
    instance.mount("/", app)
    return instance
