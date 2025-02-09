package org.example;

import javax.net.ServerSocketFactory;
import javax.net.SocketFactory;
import javax.net.ssl.SSLContext;
import javax.net.ssl.SSLServerSocketFactory;
import javax.net.ssl.SSLSocketFactory;
import java.io.*;
import java.net.InetAddress;
import java.net.ServerSocket;
import java.net.Socket;
import java.nio.charset.StandardCharsets;
import java.security.NoSuchAlgorithmException;
import java.time.Instant;
import java.util.Arrays;

public class Server {
    private static Logger log = new Logger();

    private final ServerSocketFactory serverSocketFactory;
    private final int port;

    public Server() {
        this(ServerSocketFactory.getDefault(), 7666);
    }

    public Server(ServerSocketFactory serverSocketFactory, int port) {
        this.serverSocketFactory = serverSocketFactory;
        this.port = port;
    }

    public static void main(String... args) throws NoSuchAlgorithmException {
        SSLContext sslContext = SSLContext.getDefault();
        SSLServerSocketFactory sslServerSocketFactory = (SSLServerSocketFactory) SSLServerSocketFactory.getDefault();
        System.out.println(Arrays.asList(sslServerSocketFactory.getDefaultCipherSuites()));
        System.out.println(Arrays.asList(sslServerSocketFactory.getSupportedCipherSuites()));
        if(true)return;
        new Server(SSLServerSocketFactory.getDefault(), 7666).start();
    }

    public void start() {
        log.info("Starting server");
        try (ServerSocket serverSocket = serverSocketFactory.createServerSocket(port)) {
            while (!serverSocket.isClosed()) {
                log.info("Waiting for incoming connection on %s", serverSocket.getLocalSocketAddress());
                Socket conn = serverSocket.accept();
                log.state(conn, "connection established");
                new Thread(() -> handleIncommingConnection(conn)).start();
            }
            log.info("Server stopped");
        } catch (IOException e) {
            log.error("Error:%s", e);
        }
    }

    private void handleIncommingConnection(Socket conn) {
        log.state(conn, "Connection established");
        try (BufferedReader in = new BufferedReader(new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8));
             BufferedWriter out = new BufferedWriter(new OutputStreamWriter(conn.getOutputStream(), StandardCharsets.UTF_8))) {
            for (String line = in.readLine(); line != null; line = in.readLine()) {
                log.received(conn, line);
                if ("time".equals(line)) {
                    String time = Instant.now().toString();
                    log.send(conn, time);
                    out.write(time + "\n");
                    out.flush();
                    continue;
                }

                log.send(conn, line);
                out.write(line + "\n");
                out.flush();
            }
        } catch (IOException ex) {
            log.error("Error while handling incomming connection: %s", ex);
        }
    }
}
