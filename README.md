# Chat Application

This is a simple chat application that simulates the basic functionalities of a chat app like WhatsApp. It supports real-time messaging through WebSocket and allows users to see a list of online users. Messages are styled to resemble WhatsApp, with sent messages aligned to the right and received messages to the left.

## Features

- **Real-Time Messaging**: Communicate instantly with other users using WebSocket.
- **Message Display**: Messages are styled similar to WhatsApp, with sent messages on the right and received messages on the left.
- **Username Display**: Each message displays the sender's username above the message bubble.
- **Online Users List**: A dynamic list of online users is displayed at the top.
- **Responsive Design**: The application adjusts to different screen sizes, making it usable on both desktop and mobile devices.
- **Message History**: Fetch and display previous chat history upon connection.

## Technologies Used

- **Frontend**:
  - **HTML/CSS**: For structuring and styling the chat application.
  - **JavaScript**: For client-side interactions and WebSocket communication.
  - **WebSocket**: For real-time messaging.
  - **Fetch API**: For retrieving chat history.
- **Backend**:
  - **Go**: For the WebSocket server and REST API to handle messaging and chat history.

