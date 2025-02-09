package org.example;

import java.net.Socket;
import java.time.Instant;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;

public class Logger {
    private DateTimeFormatter df;
    private Level level;

    public Logger() {
        this.df = DateTimeFormatter.ofPattern("HH:mm:ss.SSS").withZone(ZoneId.systemDefault());
        this.level = Level.INFO;
    }

    private void setFormat(DateTimeFormatter df) {
        this.df = df;
    }

    private void setLevel(Level level) {
        this.level = level;
    }

    public void trace(String msg, Object... args) {
        log(Level.INFO, msg, args);
    }

    public void debug(String msg, Object... args) {
        log(Level.INFO, msg, args);
    }

    public void info(String msg, Object... args) {
        log(Level.INFO, msg, args);
    }

    public void warn(String msg, Object... args) {
        log(Level.INFO, msg, args);
    }

    public void error(String msg, Object... args) {
        log(Level.INFO, msg, args);
    }

    public void received(Socket conn, String msg) {
        info("%s <-- %s : %s", conn.getLocalSocketAddress(), conn.getRemoteSocketAddress(), msg);
    }

    public void send(Socket conn, String msg) {
        info("%s --> %s : %s", conn.getLocalSocketAddress(), conn.getRemoteSocketAddress(), msg);
    }

    public void state(Socket conn, String msg) {
        info("%s <-> %s : %s", conn.getLocalSocketAddress(), conn.getRemoteSocketAddress(), msg);
    }

    public void log(Level level, String msg, Object... args) {
        if (level.ordinal() < this.level.ordinal()) {
            return;
        }
        Instant now = Instant.now();
        String formattedMsg = df.format(now) + " [" + Thread.currentThread().getName() + "] " + this.level.name() + " : " + String.format(msg, args);
        if (level == Level.ERROR) {
            System.err.println(formattedMsg);
        } else {
            System.out.println(formattedMsg);
        }
    }

    public enum Level {
        TRACE,
        DEBUG,
        INFO,
        WARN,
        ERROR,
    }
}
