#!/usr/bin/env python3

from flask_login import UserMixin
from datetime import datetime
from . import db

class CRUDMixin(object):
    __table_args__ = {'extend_existing': True}

    id = db.Column(db.Integer, primary_key=True)
    created_on = db.Column(db.DateTime, default=datetime.now, nullable=True)
    changed_on = db.Column(
        db.DateTime, default=datetime.now, onupdate=datetime.now, nullable=True
    )

    @classmethod
    def get_by_id(cls, id):
        if any(
            (isinstance(id, str) and id.isdigit(),
             isinstance(id, (int, float))),
        ):
            return cls.query.get(int(id))
        return None

    @classmethod
    def create(cls, **kwargs):
        instance = cls(**kwargs)
        return instance.save()

    def update(self, commit=True, **kwargs):
        for attr, value in kwargs.iteritems():
            setattr(self, attr, value)
        return commit and self.save() or self

    def save(self, commit=True):
        db.session.add(self)
        if commit:
            db.session.commit()
        return self

    def delete(self, commit=True):
        db.session.delete(self)
        return commit and db.session.commit()

class User(UserMixin, CRUDMixin, db.Model):
    __tablename__ = "users"

    email = db.Column(db.String(100), unique=True)
    password = db.Column(db.String(100))
    name = db.Column(db.String(1000))
    tasks = db.relationship("Task", backref="user")


class Task(db.Model, CRUDMixin):
    __tablename__ = "tasks"

    title = db.Column(db.String(500))
    user_id = db.Column(db.Integer, db.ForeignKey('users.id'))
    events = db.relationship("Event", backref="task")

    @property
    def event_counts(self):
        return len(self.events)

class Event(db.Model, CRUDMixin):
    __tablename__ = "events"
    task_id = db.Column(db.Integer, db.ForeignKey('tasks.id'))
