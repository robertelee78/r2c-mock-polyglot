import os

from django.core.wsgi import get_wsgi_application


def session_positive():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_session_positive"
    return get_wsgi_application()


def session_control():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_session_control"
    return get_wsgi_application()


def session_uncertain():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_session_uncertain"
    return get_wsgi_application()


def cors_positive():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_cors_positive"
    return get_wsgi_application()


def cors_control():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_cors_control"
    return get_wsgi_application()


def cors_uncertain():
    os.environ["DJANGO_SETTINGS_MODULE"] = "web.python.django_cors_uncertain"
    return get_wsgi_application()
