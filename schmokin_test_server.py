from flask import Flask
from flask import jsonify
from flask import make_response
from flask import request
from flask import redirect
import json

app = Flask(__name__)

@app.route("/simple")
def simple():
    data = {
        'status': 'UP'        
    }
    return jsonify(data)

@app.route("/pretty/redirect")
def pretty_redirect():
    data = {
        'status': 'UP'
    }
    return json.dumps(data, indent=4, sort_keys=True)


@app.route("/pretty")
def pretty():
    data = {
        'status': 'UP'
    }
    return json.dumps(data, indent=4, sort_keys=True)

@app.route("/resources/<resource>")
def resources(resource):
    return resource

@app.route("/created")
def created():
    data = {
        'status': 'UP'        
    }
    return jsonify(data), 201

@app.route("/array")
def array():
    data = [1,2,3,4,5]
    response = make_response(jsonify(data))
    response.headers.set('Content-Type', 'application/json')
    return response

@app.route('/echo', methods=['POST'])
def echo():
    data = request.get_data()
    return data

app.run(port=40000)
