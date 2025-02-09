package org.example;

import javax.net.SocketFactory;
import javax.net.ssl.SSLSocketFactory;
import java.io.*;
import java.net.Socket;
import java.nio.charset.StandardCharsets;

public class Client {

    private static final Logger log = new Logger();

    private final SocketFactory socketFactory;
    private final String host;
    private final int port;

    public Client() {
        this(SocketFactory.getDefault(), "127.0.0.1", 7666);
    }

    public Client(SocketFactory socketFactory, String host, int port) {
        this.socketFactory = socketFactory;
        this.host = host;
        this.port = port;
    }

    public static void main(String... args) throws IOException {
        System.setProperty("jdk.tls.server.protocols", "TLSv1.2");
        new Client(SSLSocketFactory.getDefault(), "localhost", 7666).start();
    }

    public void start() {
        try (Socket socket = socketFactory.createSocket(host, port)) {
            Thread consoleHandler = new Thread(() -> readConsole(socket), "ConsoleHandler");
            consoleHandler.start();

            Thread socketReceiver = new Thread(() -> readSocket(socket), "SocketReader");
            socketReceiver.start();

            log.info("Waiting for socket and console handler to finish");
            consoleHandler.join();
            socketReceiver.join();
        } catch (IOException e) {
            log.error("Error:%s", e);
        } catch (InterruptedException e) {
            log.error("Interrupted while waiting for handlers to finish");
        }
        log.info("Exiting Client");
    }

    private void readSocket(Socket conn) {
        log.state(conn, "SocketReader started");
        try (BufferedReader in = new BufferedReader(new InputStreamReader(conn.getInputStream(), StandardCharsets.UTF_8))) {
            for (String line = in.readLine(); line != null; line = in.readLine()) {
                log.received(conn, line);
            }
        } catch (IOException e) {
            log.error("Error while reading conn:%s", e);
            e.printStackTrace();
        }
        log.state(conn, "SocketReader finished");
    }

    private void readConsole(Socket conn) {
        log.state(conn, "ConsoleHandler started");
        try (BufferedReader in = new BufferedReader(new InputStreamReader(System.in, StandardCharsets.UTF_8));
             BufferedWriter out = new BufferedWriter(new OutputStreamWriter(conn.getOutputStream(), StandardCharsets.UTF_8))
        ) {
            for (String line = in.readLine(); line != null; line = in.readLine()) {
                log.send(conn, line);
                out.write(line);
                out.newLine();
                out.flush();
            }
        } catch (IOException e) {
            log.error("Error processing console:%s", e);
        }
        log.state(conn, "ConsoleHandler finished");
    }
}
