package com.microservicios.orchestrator.model;

import java.util.List;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class AlertEvent {

    private String receiver;
    private String status;
    private List<AlertItem> alerts;

    @JsonCreator
    public AlertEvent(@JsonProperty("receiver") String receiver, @JsonProperty("status") String status,
            @JsonProperty("alerts") List<AlertItem> alerts) {
        this.receiver = receiver;
        this.status = status;
        this.alerts = alerts;
    }

    // Getters y setters
    public String getReceiver() {
        return receiver;
    }

    public void setReceiver(String receiver) {
        this.receiver = receiver;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public List<AlertItem> getAlerts() {
        return alerts;
    }

    public void setAlerts(List<AlertItem> alerts) {
        this.alerts = alerts;
    }
}