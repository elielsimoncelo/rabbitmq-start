#!/usr/bin/env node

var amqp = require('amqplib/callback_api');

amqp.connect('amqp://rabbitmq:rabbitmq@localhost:5672/', function(errorConnection, connection) {
    if (errorConnection) {
        throw errorConnection;
    }

    connection.createChannel(function(errorChanner, channel) {
        if (errorChanner) {
            throw errorChanner;
        }
        
        var exchange = 'logs';

        channel.assertExchange(exchange, 'fanout', {
            durable: true,
        });

        var queue = 'file_log_queue';
        
        channel.assertQueue(queue, {
            durable: true,
        });

        channel.prefetch(1);
        channel.bindQueue(queue, exchange, 'log.file');

        console.log(" [*] Waiting for messages in %s. To exit press CTRL+C", queue);
        
        channel.consume(queue, function(msg) {
            console.log(" [x] Received %s", msg.content.toString());

            setTimeout(function() {
                console.log(" [x] Done\n");
                channel.ack(msg);
            }, 5 * 1000);
        }, {
            noAck: false
        });
    });
});