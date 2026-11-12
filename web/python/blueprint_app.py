import os

from flask import Blueprint, Flask, Response, abort, redirect, request, session
from flask_cors import CORS
from flask_principal import Permission, Principal, RoleNeed
from flask_wtf.csrf import CSRFProtect
from markupsafe import Markup, escape
from sqlalchemy import Boolean, String, create_engine
from sqlalchemy.orm import DeclarativeBase, Mapped, Session, mapped_column

from .shared import deny_policy, object_policy, origin_policy, secure_policy, text_policy


class Base(DeclarativeBase):
    pass


class Account(Base):
    __tablename__ = "blueprint_accounts"

    id: Mapped[int] = mapped_column(primary_key=True)
    display: Mapped[str] = mapped_column(String(120))
    is_admin: Mapped[bool] = mapped_column(Boolean, default=False)


engine = create_engine(os.environ["O84_DATABASE_URL"])

app = Flask(__name__)
app.config["SECRET_KEY"] = os.environ["O84_SESSION_SECRET"]
csrf = CSRFProtect(app)
principal = Principal(app)
editor_permission = Permission(RoleNeed("editor"))
policy_permission = Permission(RoleNeed("editor"))
bp = Blueprint("accounts", __name__, url_prefix="/accounts")


@bp.get("/html/positive")
def html_positive():
    return Markup(request.args["html"])


@bp.get("/html/control")
def html_control():
    return Markup(escape(request.args["html"]))


@bp.get("/html/uncertain")
def html_uncertain():
    return Markup(text_policy(request.args["html"]))


@bp.post("/csrf/positive")
@csrf.exempt
def csrf_positive():
    session["last_action"] = request.form["action"]
    return "saved"


@bp.post("/csrf/control")
def csrf_control():
    session["last_action"] = request.form["action"]
    return "saved"


@bp.route("/csrf/uncertain", methods=["GET", "POST"])
def csrf_uncertain():
    session["last_action"] = request.form["action"]
    return "saved"


@bp.get("/redirect/positive")
def redirect_positive():
    return redirect(request.args["next"])


@bp.get("/redirect/control")
def redirect_control():
    return redirect("/home")


@bp.get("/redirect/uncertain")
def redirect_uncertain():
    return redirect(text_policy(request.args["next"]))


@bp.post("/authorization/positive")
def authorization_positive():
    if not editor_permission.can():
        Response("forbidden", status=403)
    session["changed"] = request.form["value"]
    return "updated"


@bp.post("/authorization/control")
def authorization_control():
    if not editor_permission.can():
        abort(403)
    session["changed"] = request.form["value"]
    return "updated"


@bp.post("/authorization/uncertain")
def authorization_uncertain():
    if not policy_permission.can():
        deny_policy(policy_permission)
    session["changed"] = request.form["value"]
    return "updated"


@bp.post("/binding/positive")
def binding_positive():
    account = Account()
    account.is_admin = False
    account.__dict__.update(request.get_json())
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@bp.post("/binding/control")
def binding_control():
    account = Account()
    account.is_admin = False
    account.display = request.get_json()["display"]
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


@bp.post("/binding/uncertain")
def binding_uncertain():
    account = Account()
    account.is_admin = False
    account.__dict__.update(object_policy(request.get_json()))
    db = Session(engine)
    db.add(account)
    db.commit()
    return {"saved": True}


app.register_blueprint(bp)


def session_positive():
    instance = Flask(__name__)
    instance.config.update(SESSION_COOKIE_SECURE=False)
    return instance


def session_control():
    instance = Flask(__name__)
    instance.config.update(SESSION_COOKIE_SECURE=True, SESSION_COOKIE_HTTPONLY=True)
    return instance


def session_uncertain():
    instance = Flask(__name__)
    instance.config.update(SESSION_COOKIE_SECURE=False if not secure_policy() else True)
    return instance


def cors_positive():
    instance = Flask(__name__)
    CORS(instance, origins="*", supports_credentials=True)
    return instance


def cors_control():
    instance = Flask(__name__)
    CORS(instance, origins=["https://trusted.example"], supports_credentials=True)
    return instance


def cors_uncertain():
    instance = Flask(__name__)
    CORS(instance, origins=origin_policy(), supports_credentials=True)
    return instance
