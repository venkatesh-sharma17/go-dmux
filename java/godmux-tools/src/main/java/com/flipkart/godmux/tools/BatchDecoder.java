package com.flipkart.godmux.tools;

import java.nio.ByteBuffer;
import java.util.Arrays;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public enum BatchDecoder {
    DECODE;

    private static final Logger logger = LoggerFactory.getLogger(BatchDecoder.class);

    public byte[][] decode(byte[] payload) {
        logger.debug("got payload len = {} ", payload.length);
        int intSize = 4;
        int size = ByteBuffer.wrap(payload, 0, intSize).getInt();
        logger.debug("got batch size {}", size);
        int[] indvidualSizes = new int[size];
        byte[][] output = new byte[size][];
        int offset = intSize;
        for (int i = 0; i < size; i++) {
            int len = ByteBuffer.wrap(payload, offset, intSize).getInt();
            indvidualSizes[i] = len;
            offset += intSize;
            byte[] data = new byte[len];
            System.arraycopy(payload, offset, data, 0, len);
            output[i] = data;
            offset += len;
        }
        logger.debug("got individual payloads of sizes {}", Arrays.toString(indvidualSizes));
        return output;
    }
}
