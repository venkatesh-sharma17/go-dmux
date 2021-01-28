package com.flipkart.godmux.tools.request;

import java.util.List;

public class BatchRequest {

    private int partition;
    private List<DmuxRequest> dmuxRequestList;

    public BatchRequest(int partition, List<DmuxRequest> dmuxRequestList) {
        this.partition = partition;
        this.dmuxRequestList = dmuxRequestList;
    }

    public int getPartition() {
        return partition;
    }

    public List<DmuxRequest> getDmuxRequestList() {
        return dmuxRequestList;
    }
}
