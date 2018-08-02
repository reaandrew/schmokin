from flask import Flask
from flask import jsonify
from flask import make_response

app = Flask(__name__)

@app.route("/simple")
def simple():
    data = {
        'status': 'UP'        
    }
    return jsonify(data)

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

app.run(port=40000)
