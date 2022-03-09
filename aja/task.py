#!/usr/bin/env python3
from flask import Blueprint, render_template, request, redirect
from . import db
from .models import Task, Event
from flask_login import login_required, current_user
import logging
from datetime import datetime

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

    # Checking if a "float" event already exists for the task for the given day.

    # Retrieving the events
    task_events = Event.query.filter_by(task_id=task.id).all()
    logging.warn(task_events)

    # Looking for one with the timestamp within the last day.
    for event_elem in task_events:
        if event_elem.created_on.date() == datetime.today().date():
            logging.warn("found event!")
            curr_event = event_elem
            break

    try:
        curr_event
    except NameError:
        curr_event = Event(task=task)
        curr_event.event_duration = 0
    else:
        logging.warn(curr_event.event_duration)

    curr_event.event_duration = curr_event.event_duration  + 0.5
    curr_event.save()
    return redirect('/tasks')
