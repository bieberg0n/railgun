import socket
from urllib.parse import urlparse#, urlunparse
from threading import Thread
from multiprocessing import Process
import re
import requests
import ssl


def parse_header(raw_headers):
	request_lines = raw_headers.split('\r\n')
	first_line = request_lines[0].split(' ')
	method = first_line[0]
	full_path = first_line[1]
	version = first_line[2]
	print("%s %s" % (method, full_path))
	(scm, netloc, path, params, query, fragment) \
		= urlparse(full_path, 'http')
	if method == 'CONNECT':
		address = (path.split(':')[0], int(path.split(':')[1]))
	else:
		# 如果url中有‘：’就指定端口，没有则为默认80端口
		i = netloc.find(':')
		if i >= 0:
			address = netloc[:i], int(netloc[i + 1:])
		else:
			address = netloc, 80
	return method, version, scm, address, path, params, query, fragment


host_p = re.compile('http://.+?/')
def get_rawheader_met(conn):
	headers = ''
	while 1:
		buf = conn.recv(2048).decode('utf-8')
		headers += buf
		# if headers[-4:] == '\r\n\r\n':
		if len(buf) < 2048:
			break

	method, version, scm, address, path, params, query, fragment =\
		parse_header(headers)
	return headers, method


def get_headers(raw_headers):
	headers = raw_headers.replace(
		'Proxy-Connection: keep-alive', 'Connection: close'
		).replace(
			'keep-alive', 'close'
			)
	headers = host_p.sub('/', headers)
	return headers


s = requests.session()
context = ssl.SSLContext(ssl.PROTOCOL_TLSv1)
# context.load_cert_chain(certfile="cacert.pem", keyfile="CA.crt")
context.load_cert_chain(certfile="server1.pem", keyfile="server.key")
def get_resp_for_ssl(connect):
	# while 1:
	conn = context.wrap_socket(connect, server_side = True)
	req = b''
	buf = 1
	while buf:
		buf = conn.read()
		print(buf)
		req += buf
		if req.endswith(b'\r\n\r\n'):
			req = get_headers(req.decode('utf-8'))
			print(req)
			r = s.post('http://127.0.0.1:8080', data={'req':req,'ssl':1})
			conn.send(r.content)
			break
	conn.shutdown(socket.SHUT_RDWR)
	conn.close()


def handle_connection(conn):
	# 从socket读取头
	# print(conn)
	raw_headers, method = get_rawheader_met(conn)
	# print(repr(req_headers))
	if method == 'CONNECT':
		conn.sendall(b'HTTP/1.1 200 Connection established\r\n\r\n')
		print('发给浏览器连接确认')
		get_resp_for_ssl(conn)

	else:
		headers = get_headers(raw_headers)
		r = s.post('http://127.0.0.1:8080', data={'req':headers})
		# r = s.post('https://mc-bieber.rhcloud.com', data={'req':req_headers})
		# print(r.content)
		# r = s.post('http://127.0.0.1:8080', data={'req':req})
		# return r.content
		conn.sendall(r.content)
		conn.close()
	exit()


	
	# print(address, req_headers)
	soc = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	# soc.settimeout(1)
	# 尝试连接
	try:
		# print('connect {}'.format(address))
		soc.connect(address)
		# print('connect scuess')
	except socket.error:
		# conn.sendall("HTTP/1.1" + str(arg[0]) + " Fail\r\n\r\n")
		conn.close()
		soc.close()
	else:  # 若连接成功
		def serv(conn, soc):
			while 1:
				data = conn.recv(1024*8)
				print('cli:', len(data))
				if not data:
					break
				soc.sendall(data)
		def cli(conn, soc):
			while 1:
				data = soc.recv(1024*8)
				print('serv:', len(data))
				if not data:
					break
				conn.sendall(data)
				
		if method == 'CONNECT':
			conn.sendall(b'HTTP/1.1 200 Connection established\r\n\r\n')
			# try:
			Process(target=serv, args=(conn, soc,)).start()
			Process(target=cli, args=(conn, soc,)).start()
			# while 1:
				# data = conn.recv(99999)
				# if not data:
				# 	break
				# soc.sendall(data)
				# data  = soc.recv(99999)
				# conn.sendall(data)
			# except:
			# 	conn.close()
			# 	soc.close()
					
		else:
			soc.sendall(req_headers.encode('utf-8'))
			# 发送完毕, 接下来从soc读取服务器的回复
			# 建立个缓冲区
			# data = b''
			# while 1:
			# try:
					# buf = soc.recv(1024*8)
					# if not buf:
					# 	break
					# conn.sendall(buf)
			cli(conn, soc)
			# except socket.error:
				# break
				# pass
				# data += buf
			# 转发给客户端
			print(address, 'ok data')
			
		soc.close()
		conn.close()


def server():
	s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
	s.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
	s.bind(('127.0.0.1', 8000))
	s.listen(1500)
	print("Serving at 127.0.0.1:8000")
	while 1:
		try:
			conn, addr = s.accept()
			# print(addr)
			Thread(target=handle_connection,args=(conn,)).start()
			# handle_connection(conn)
		except KeyboardInterrupt:
			s.close()
			print("Bye...")
			break

if __name__ == '__main__':
	server()
