package com.flipkart.godmux.tools.request;

public class DmuxRequest {

    private String key;
    private long offset;
    private byte[] data;

    public DmuxRequest(String key, long offset, byte[] data) {
        this.key = key;
        this.offset = offset;
        this.data = data;
    }

    public String getKey() {
        return key;
    }

    public long getOffset() {
        return offset;
    }

    public byte[] getData() {
        return data;
    }
}
