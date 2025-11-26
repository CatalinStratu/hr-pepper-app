package com.hr.pepperinterview.models;

import com.google.gson.annotations.SerializedName;

public class Question {
    @SerializedName("id")
    private String id;

    @SerializedName("text")
    private String text;

    @SerializedName("category")
    private String category;

    @SerializedName("difficulty")
    private String difficulty;

    @SerializedName("time_limit")
    private int timeLimit;

    @SerializedName("order")
    private int order;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    public String getCategory() {
        return category;
    }

    public void setCategory(String category) {
        this.category = category;
    }

    public String getDifficulty() {
        return difficulty;
    }

    public void setDifficulty(String difficulty) {
        this.difficulty = difficulty;
    }

    public int getTimeLimit() {
        return timeLimit;
    }

    public void setTimeLimit(int timeLimit) {
        this.timeLimit = timeLimit;
    }

    public int getOrder() {
        return order;
    }

    public void setOrder(int order) {
        this.order = order;
    }
}
