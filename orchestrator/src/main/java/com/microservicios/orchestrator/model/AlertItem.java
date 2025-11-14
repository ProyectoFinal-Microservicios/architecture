package com.microservicios.orchestrator.model;

import java.util.Map;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;

public class AlertItem {

    private String status;
    private Map<String, String> labels;
    private Map<String, String> annotations;

    @JsonCreator
    public AlertItem(@JsonProperty("status") String status,
                    @JsonProperty("labels") Map<String, String> labels,
                    @JsonProperty("annotations") Map<String, String> annotations) {
        this.status = status;
        this.labels = labels;
        this.annotations = annotations;
    }

    // Getters y setters
    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

    public Map<String, String> getLabels() {
        return labels;
    }

    public void setLabels(Map<String, String> labels) {
        this.labels = labels;
    }

    public Map<String, String> getAnnotations() {
        return annotations;
    }

    public void setAnnotations(Map<String, String> annotations) {
        this.annotations = annotations;
    }
}
