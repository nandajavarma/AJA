#!/usr/bin/env python3
from flask import Blueprint, render_template, request
from . import db
from flask_login import login_required, current_user
import logging

main = Blueprint('main', __name__)

@main.route('/')
def index():
    return app.send_static_file("index.html")

@main.route('/profile')
@login_required
def profile():
    return app.send_static_file("profile.html", name=current_user.name)

@main.route('/hello')
def hello():
    return {"result": "Hello world"}
