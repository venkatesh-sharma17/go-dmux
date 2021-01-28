package com.flipkart.godmux.tools;

import com.flipkart.godmux.tools.request.BatchRequest;
import com.flipkart.godmux.tools.request.DmuxRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.nio.ByteBuffer;
import java.util.ArrayList;
import java.util.List;

public class DmuxDecoder {

    public DmuxDecoder() {
        resetOffset();
    }

    private static final Logger logger = LoggerFactory.getLogger(DmuxDecoder.class);
    private int currentOffset = 0;

    private void resetOffset() {
        this.currentOffset = 0;
    }

    private DmuxRequest decode(byte[] payload, boolean keepOffset) {
        if (!keepOffset)
            resetOffset();
        return decode(payload);
    }

    public DmuxRequest decode(byte[] payload) {
        int len = readInteger(payload);
        long zkOffset = readLong(payload);
        int keySize = readInteger(payload);
        String key = new String(readByteArray(payload, keySize));
        byte[] data = readByteArray(payload, len);
        logger.debug("got individual payload of sizes {}", len);
        return new DmuxRequest(key, zkOffset, data);
    }

    public BatchRequest batchDecode(byte[] payload) {
        logger.debug("got payload len = {} ", payload.length);
        resetOffset();
        int partition = readInteger(payload);
        int size = readInteger(payload);
        logger.debug("got batch size {}", size);
        List<DmuxRequest> dmuxRequestList = new ArrayList<>();
        for (int i = 0; i < size; i++) {
            dmuxRequestList.add(decode(payload, true));
        }
        return new BatchRequest(partition, dmuxRequestList);
    }

    private int readInteger(byte[] payload) {
        int val = ByteBuffer.wrap(payload, currentOffset, 4).getInt();
        currentOffset += 4;
        return val;
    }

    private long readLong(byte[] payload) {
        long val = ByteBuffer.wrap(payload, currentOffset, 8).getLong();
        currentOffset += 8;
        return val;
    }

    private byte[] readByteArray(byte[] payload, int length) {
        byte[] data = new byte[length];
        System.arraycopy(payload, currentOffset, data, 0, length);
        currentOffset += length;
        return data;
    }

}
