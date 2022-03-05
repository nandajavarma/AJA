#!/usr/bin/env python3
from flask import Blueprint, render_template, request, redirect
from . import db
from .models import Task, Event
from flask_login import login_required, current_user
import logging

task = Blueprint('task', __name__)

@task.route('/tasks')
@login_required
def index():
    tasks = Task.query.filter_by(user=current_user).all()
    logging.warn(tasks)
    return render_template("tasks.html", tasks=tasks)

@task.route('/tasks', methods=["POST"])
@login_required
def handle_add_task():
    task_input = request.form.get("text")
    logging.info("Creating new task")
    new_task = Task(title=task_input, user=current_user)
    new_task.save()
    return index()

@task.route('/event', methods=["POST"])
@login_required
def handle_add_event():
    task_id = request.form.get("text")
    task = Task.get_by_id(int(task_id))
    new_event = Event(task=task)
    new_event.save()
    return redirect('/tasks')
