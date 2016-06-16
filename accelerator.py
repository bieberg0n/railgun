from flask import Flask, request
from gevent.wsgi import WSGIServer
import re
# import socket
from gevent import ssl, socket#,monkey
from socket import getaddrinfo
# monkey.patch_socket
import time

app = Flask(__name__)

path_p = re.compile('^.+? (.+?) ')
host_p = re.compile('Host: (.+?)\r')
def get_host_port(req, https=False):
	host = host_p.findall(req)[0]
	if ':' in host:
		host_without_port = host.split(':')[0]
		port = int(host.split(':')[1])
		return host_without_port, port
	else:
		port = 443 if https else 80
		return host, port


def get_resp(req, https=False):
	url = '{host}{path}'.format(
		host = req.split('\r\n')[1].split(' ')[1],
		path = req.split('\r\n')[0].split(' ')[1])
	info = '{method} {url}'.format(
		method = req.split(' ')[0], url = url)
	# print(info)
	host, port = get_host_port(req, https)
	if https:
		s = ssl.wrap_socket(socket.socket(),
							ssl_version=ssl.PROTOCOL_TLSv1)
		s.settimeout(5)
	else:
		s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	try:
		s.connect((host, port))
		# print('connect ok')
	except socket.timeout as e:
			print(e)
			return
	s.sendall(req.encode('utf-8'))
	# print('send ok')
	
	resp = b''
	while 1:
		# buf and not buf.startswith(b'WebSocket')\
		  # and not buf.endswith(b'\r\n\r\n'):
		try:
			buf = s.recv(1024*8)
		except socket.timeout as e:
			print(e)
			break
		
		# print('{}......{}'.format(buf[:50],buf[-50:]))
		resp += buf
		if not buf or\
		   buf.startswith(b'WebSocket') and buf.endswith(b'\r\n\r\n'):
			break
	# print('RECV {}'.format(url))
	return resp, info


def getip(name):
	try:
		ip = getaddrinfo(name,None)[0][-1][0]
	except socket.gaierror:
		# pass
		ip = '127.0.0.1'
	return ip


@app.route('/',methods=['GET','POST'])
def index():
	if request.method == 'GET':
		return 'Welcome to bjdns!'
	elif request.method == 'POST':
		if request.form.get('n'):
			return getip(request.form['n'])
		else:
			return 'None'


@app.route('/proxy', methods=['GET', 'POST'])
def accelerator():
	if request.method == 'GET':
		return 'Railgun.'
	elif request.method == 'POST':
		if request.form.get('ssl'):
			# print('ssl')
			resp, info= get_resp(request.form['req'], https=True)
		else:
			resp, info= get_resp(request.form['req'])

		print(request.remote_addr,
			  '- - [{}]'.format(time.strftime('%Y-%m-%d %H:%M:%S')),
			  '"{}"'.format(info) )
		return resp


# http_server = WSGIServer(('', 8080), app)
http_server = WSGIServer(('127.12.99.1', 8080), app)
http_server.serve_forever()
# app.run(host='127.0.0.1', port=8080, debug=True)
