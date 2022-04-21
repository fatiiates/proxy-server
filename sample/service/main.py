#!/usr/bin/env python
# encoding: utf-8
import json
from flask import Flask, request

app = Flask(__name__)

@app.route('/sum', methods=['POST'])
def sum():
    args = request.args
    return json.dumps({'res': int(args.get('n1')) + int(args.get('n2'))})

app.run(host="0.0.0.0")