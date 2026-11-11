import importlib
import os


def _configured(variable):
    module, _, member = os.environ[variable].partition(":")
    return getattr(importlib.import_module(module), member)


def text_policy(value):
    return _configured("O84_TEXT_POLICY")(value)


def object_policy(values):
    return _configured("O84_OBJECT_POLICY")(values)


def origin_policy():
    return _configured("O84_ORIGIN_POLICY")()


def secure_policy():
    return _configured("O84_SECURE_POLICY")()


def scope_policy(scope):
    return _configured("O84_SCOPE_POLICY")(scope)


def token_scopes(token):
    return _configured("O84_TOKEN_SCOPES")(token)


def field_policy(model):
    return _configured("O84_FIELD_POLICY")(model)


def csrf_policy(session, headers):
    return _configured("O84_CSRF_POLICY")(session, headers)


def user_policy(user):
    return _configured("O84_USER_POLICY")(user)


def deny_policy(subject):
    return _configured("O84_DENY_POLICY")(subject)


def configured_loader(handle):
    return _configured("O84_CONFIG_LOADER")(handle)
