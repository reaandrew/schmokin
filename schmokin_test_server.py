from flask import Flask
from flask import jsonify
app = Flask(__name__)

@app.route("/simple")
def simple():
    data = {
        'status': 'UP'        
    }
    return jsonify(data)

@app.route("/array")
def array():
    data = [1,2,3,4,5]
    return jsonify(data)

app.run(port=40000)
