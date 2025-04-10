
package com.example.testapp;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.io.PrintWriter;
import java.net.Socket;

public class TestApp {
    public static void main(String[] args) {
        String host = "localhost";
        int port = 8000;

        try (Socket socket = new Socket(host, port)) {
            OutputStream outputStream = socket.getOutputStream();
            PrintWriter writer = new PrintWriter(outputStream, true);

            writer.println("GET / HTTP/1.1");
            writer.println("Host: " + host);
            writer.println();

            writer.flush();

            InputStream inputStream = socket.getInputStream();

            byte[] buffer = new byte[1024];
            int bytesRead;
            while ((bytesRead = inputStream.read(buffer)) != -1) {
                System.out.write(buffer, 0, bytesRead);
            }

            inputStream.close();

        } catch (IOException e) {
            e.printStackTrace();
        }

    }
}
