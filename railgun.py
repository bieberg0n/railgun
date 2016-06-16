# import socket
from urllib.parse import urlparse#, urlunparse
# from threading import Thread
# from multiprocessing import Process
import re
import os
import requests
# import ssl
from certutil import CertUtil
import time
# import queue
from gevent import ssl, monkey
from gevent.server import StreamServer
monkey.patch_socket()

def parse_header(raw_headers):
	request_lines = raw_headers.split('\r\n')
	first_line = request_lines[0].split(' ')
	method = first_line[0]
	if len(first_line) >= 2: 
		full_path = first_line[1]
	else:
		full_path = first_line[0]
	try:
		host = request_lines[1].split(' ')[1]
	except IndexError:
		# print(raw_headers)
		host = ''
	# print(host)
	# version = first_line[2]
	# print("%s %s" % (method, full_path))
	(scm, netloc, path, params, query, fragment) \
		= urlparse(full_path, 'http')
	# if method == 'CONNECT':
	# 	address = (path.split(':')[0], int(path.split(':')[1]))
	# else:
	# 	# 如果url中有‘：’就指定端口，没有则为默认80端口
	# 	i = netloc.find(':')
	# 	if i >= 0:
	# 		address = netloc[:i], int(netloc[i + 1:])
	# 	else:
	# 		address = netloc, 80
	return method, full_path, host#, version, scm, address, path, params, query, fragment


def get_rawheader_met(conn):
	headers = ''
	while 1:
		buf = conn.recv(2048).decode('utf-8')
		headers += buf
		# if headers[-4:] == '\r\n\r\n':
		if len(buf) < 2048:
			break

	method, _, _ = parse_header(headers)#, version, scm, address, path, params, query, fragment =\
	return headers, method


host_p = re.compile('http://.+?/')
connection_p = re.compile('Connection: .+?\r')
def get_headers(raw_headers):
	headers = raw_headers.replace(
		'Proxy-Connection: keep-alive', 'Connection: close'
		)# .replace(
		 # 	'keep-alive', 'close'
		 # 	)
	headers = connection_p.sub('Connection: Close\r', headers)
	headers = headers.split('\n')
	headers[0] = host_p.sub('/', headers[0])
	return '\n'.join(headers)


s = requests.session()
# context.load_cert_chain(certfile="cacert.pem", keyfile="CA.crt")
commonname_p = re.compile('Host: (.+?)\r')
def get_req_for_ssl(connect, raw_headers):
	# while 1:
	host = commonname_p.findall(raw_headers)[0].replace(':443','')
	# print(host)
	# context.load_cert_chain(certfile=".zhihu.com.pem")#, keyfile="go.key")
	# file = CertUtil.get_cert(host)
	# print(file)
	context = ssl.SSLContext(ssl.PROTOCOL_TLSv1)
	# context.load_cert_chain(certfile='certs/.baidu.com.crt')
	context.load_cert_chain(certfile=CertUtil.get_cert(host))
	conn = context.wrap_socket(connect, server_side = True)
	req = b''
	buf = 1
	while buf:
		buf = conn.recv(1024*8)
		# print(buf)
		req += buf
		if not req.startswith(b'POST') and b'\r\n\r\n' in req\
		   or req.startswith(b'POST') and 1< len(buf) < 2048\
		   and not buf.endswith(b'\r\n\r\n'):
			break
		# if 1< len(buf) < 2048:
		# if not buf:
	try:
		req = get_headers(req.decode('utf-8'))
	except UnicodeDecodeError:
		# break
		return None, conn
	else:
		# if 'apilive' in req:
		if 'websocket' in req:
			# print(req)
			return None, conn
		else:
			return req, conn
	# print(req)
	# conn.sendall(r.content)
	# break
	# conn.close()


def handle_connection(conn, addr):
	# while 1:
		# conn, addr = queue.get()
	raw_headers, method = get_rawheader_met(conn)
	if raw_headers:
		pass
	else:
		conn.close()
		# continue
		return
	if method == 'CONNECT':
		conn.sendall(b'HTTP/1.1 200 Connection established\r\n\r\n')
		req, conn = get_req_for_ssl(conn, raw_headers)
		# r = s.post('http://127.0.0.1:8080',
		# r = s.post('https://mc-bieber.rhcloud.com/',
		if req:
			try:
				r = s.post('http://rss.bjgong.tk/proxy', data={'req':req,'ssl':1}, stream=True)
				# r = s.post('http://bjgong.tk:8080/proxy', data={'req':req,'ssl':1})
				# r = s.post('http://127.0.0.1:8080/proxy', data={'req':req,'ssl':1})
			except requests.exceptions.ConnectionError:
				conn.close()
				return
		else:
			# conn.shutdown(socket.SHUT_RDWR)
			conn.close()
			# continue
			return
			# r = None

	else:
		req = get_headers(raw_headers)
		# r = s.post('http://127.0.0.1:8080', data={'req':headers})
		try:
			r = s.post('http://rss.bjgong.tk/proxy', data={'req':req}, stream=True)
			# r = s.post('http://bjgong.tk:8080/proxy', data={'req':req})
			# r = s.post('http://127.0.0.1:8080/proxy', data={'req':req})
		except requests.exceptions.ConnectionError:
			conn.close()
			return

		# r = s.post('https://mc-bieber.rhcloud.com', data={'req':headers})
		# print(r.content)
		# r = s.post('http://127.0.0.1:8080', data={'req':req})
		# return r.content
	# if r:
	for resp in r.iter_content(chunk_size=1024*8):
		conn.sendall(resp)
	# else:
	# 	pass
	method, full_path, host = parse_header(req)
	print( addr[0],
		   '[{}]'.format(time.strftime('%Y-%m-%d %H:%M:%S')),
		   '"{} {}{}"'.format(method, host, full_path) )#,
		   # '({})'.format(i) )
	# conn.shutdown(socket.SHUT_RDWR)
	conn.close()


# def server():
# 	s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
# 	s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
# 	s.bind(('0.0.0.0', 8087))
# 	s.listen(1500)
# 	print("Serving at 0.0.0.0:8087")
# 	# queue_list = [ queue.Queue() for i in range(6) ]
# 	# for i in range(6):
# 	# 	t = Thread(target=handle_connection,args=(queue_list[i],
# 	# 											  i,))
# 	# 	t.setDaemon(True)
# 	# 	t.start()

# 	while 1:
# 		# for q in queue_list:
# 		try:
# 			conn, addr = s.accept()
# 			# q.put((conn,addr))
# 			# print(addr)
# 			Thread(target=handle_connection,args=(conn,addr,)).start()
# 			# handle_connection(conn)
# 		except KeyboardInterrupt:
# 			s.close()
# 			print("Bye...")
# 			# break
# 			return

if __name__ == '__main__':
	# server()
	if os.path.isfile('CA.crt'):
		pass
	else:
		CertUtil.dump_ca()
	if os.path.isdir('certs'):
		pass
	else:
		os.mkdir('certs')
	StreamServer( ('0.0.0.0', 8087), handle_connection ).serve_forever()
