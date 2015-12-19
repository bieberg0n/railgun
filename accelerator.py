from flask import Flask, request
from gevent.wsgi import WSGIServer
import re
# import requests
import socket
import ssl

app = Flask(__name__)

path_p = re.compile('^.+? (.+?) ')
host_p = re.compile('Host: (.+?)\r')
def get_headers(req):
	# print(req)
	headers_raw = req.split('\r\n\r\n')[0].split('\r\n')
	method = headers_raw[0].split(' ')[0]
	headers = { i.split(': ')[0]:i.split(': ')[1]
				for i in headers_raw[2:]
				if i != ''# and 'GET' not in i and 'Host' not in i
				# and 'POST' not in i and 'CONNECT' not in i
				}
	url = host_p.findall(req)[0] + path_p.findall(req)[0]
	if method == 'POST':
		data = req.split('\r\n\r\n')[1].split('&')
		data = { i.split('=')[0]:i.split('=')[1] for i in data }
		print(url, headers, data)
		return method, url, headers, data
	else:
		print(method, url)
		return method, url, headers, ''


def get_host_port(req, https=False):
	host = host_p.findall(req)[0]
	if ':' in host:
		host_without_port = host.split(':')[0]
		port = int(host.split(':')[1])
		return host_without_port, port
	else:
		port = 443 if https else 80
		return host, port


# def get_resp_with_ssl(req):
# 	print(req)
# 	host, port = get_host_port(req)
# 	# s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# 	s = ssl.wrap_socket(socket.socket())
# 	s.connect((host, port))
# 	print('ssl contfinish')
# 	s.sendall(req.encode('utf-8'))
	
# 	resp = b''
# 	buf = 1
# 	while buf:
# 		buf = s.recv(1024*1024*8)
# 		resp += buf
# 	return resp


def get_resp(req, https=False):
	# req = request.form['req']
	print(req)
	host, port = get_host_port(req, https)
	if https:
		s = ssl.wrap_socket(socket.socket())
	else:
		s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	s.connect((host, port))
	s.sendall(req.encode('utf-8'))
	print('ssl connect finish')
	
	resp = b''
	buf = 1
	while buf:
		buf = s.recv(1024*1024*8)
		print(buf[-50:])
		resp += buf
	return resp


# s = requests.session()
# resp = [ 'Server', 'Date', 'Content-Type', 'Vary', 'Expires',
# 		'cache-control' ]
@app.route('/', methods=['GET', 'POST'])
def accelerator():
	if request.method == 'GET':
		return 'Railgun.'
	elif request.method == 'POST':
		if request.form.get('ssl'):
			print('ssl')
			# print(request.form['req'])
			return get_resp(request.form['req'], https=True)
		else:
			return get_resp(request.form['req'])
	# elif request.method == 'CONNECT':
		# return

		
		method, url, headers, data = get_headers(request.form['req'])
		if method == 'GET':
			r = s.get('http://{}'.format(url))
		elif method == 'POST':
			r = s.post('http://{}'.format(url), data=data)
			# print(r.text)
		ok = 'OK' if r.ok else ''
		if r.headers.get('content-length'):
			del r.headers['content-length']
		# del r.headers['expires']
		# del r.headers['x-cache']
		# del r.headers['x-ec-custom-error']
		# del r.headers['vary']
		# del r.headers['last-modified']
		# del r.headers['etag']
		# del r.headers['cache-control']
		if r.headers.get('transfer-encoding'):
			del r.headers['transfer-encoding']
		# del r.headers['date']
		# del r.headers['content-type']
		if r.headers.get('content-encoding'):
			del r.headers['content-encoding']
		headers = [ '{}: {}'.format(i,r.headers[i]) for i in r.headers ]
		# headers = ['{}: {}'.format(
		# 	'content-length',r.headers['content-length']),
		# 	'{}: {}'.format('expires',r.headers['expires'])
		# 	]
		headers = '\r\n'.join(headers) + '\r\n\r\n'
		res = 'HTTP/1.1 {code} {ok}\r\n{headers}'.format(code=r.status_code,ok=ok,headers=headers)
		# res = 'HTTP/1.1 {code} {ok}\r\n\r\n'.format(code=r.status_code,ok=ok)
		res = res.encode('utf-8') + r.content
		# print(res)
		return res
		# return r.content


# http_server = WSGIServer(('', 8080), app)
# http_server.serve_forever()
app.run(host='127.0.0.1', port=8080, debug=True)
