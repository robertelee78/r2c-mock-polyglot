import json

from flask import Flask

from .shared import configured_loader


def json_positive():
    app = Flask(__name__)
    app.config.from_file("settings-positive.json", load=json.load)
    return app


def json_control():
    app = Flask(__name__)
    app.config.from_file("settings-control.json", load=json.load)
    return app


def json_uncertain():
    app = Flask(__name__)
    app.config.from_file("settings-uncertain.json", load=configured_loader)
    return app
