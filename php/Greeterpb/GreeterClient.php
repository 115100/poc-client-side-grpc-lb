<?php
// GENERATED CODE -- DO NOT EDIT!

namespace Greeterpb;

/**
 */
class GreeterClient extends \Grpc\BaseStub {

    /**
     * @param string $hostname hostname
     * @param array $opts channel options
     * @param \Grpc\Channel $channel (optional) re-use channel object
     */
    public function __construct($hostname, $opts, $channel = null) {
        parent::__construct($hostname, $opts, $channel);
    }

    /**
     * @param \Greeterpb\GreetRequest $argument input argument
     * @param array $metadata metadata
     * @param array $options call options
     */
    public function Greet(\Greeterpb\GreetRequest $argument,
      $metadata = [], $options = []) {
        return $this->_simpleRequest('/greeterpb.Greeter/Greet',
        $argument,
        ['\Greeterpb\GreetReply', 'decode'],
        $metadata, $options);
    }

}
