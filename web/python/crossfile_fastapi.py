from fastapi import FastAPI
from sqlalchemy.orm import Session

from .crossfile_asgi_models import Account, AccountPayload, engine
from .shared import object_policy

app = FastAPI()


@app.post("/crossfile/binding/positive")
def binding_positive(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.__dict__.update(payload.model_dump())
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@app.post("/crossfile/binding/control")
def binding_control(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.display = payload.display
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@app.post("/crossfile/binding/uncertain")
def binding_uncertain(payload: AccountPayload):
    account = Account()
    account.is_admin = False
    account.__dict__.update(object_policy(payload.model_dump()))
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}
