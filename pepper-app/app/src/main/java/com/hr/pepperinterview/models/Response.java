package com.hr.pepperinterview.models;

import com.google.gson.annotations.SerializedName;

public class Response {
    @SerializedName("id")
    private String id;

    @SerializedName("interview_id")
    private String interviewId;

    @SerializedName("question_id")
    private String questionId;

    @SerializedName("response_text")
    private String responseText;

    @SerializedName("duration")
    private int duration;

    @SerializedName("sentiment_score")
    private Double sentimentScore;

    @SerializedName("confidence_score")
    private Double confidenceScore;

    @SerializedName("analysis")
    private String analysis;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getInterviewId() {
        return interviewId;
    }

    public void setInterviewId(String interviewId) {
        this.interviewId = interviewId;
    }

    public String getQuestionId() {
        return questionId;
    }

    public void setQuestionId(String questionId) {
        this.questionId = questionId;
    }

    public String getResponseText() {
        return responseText;
    }

    public void setResponseText(String responseText) {
        this.responseText = responseText;
    }

    public int getDuration() {
        return duration;
    }

    public void setDuration(int duration) {
        this.duration = duration;
    }

    public Double getSentimentScore() {
        return sentimentScore;
    }

    public void setSentimentScore(Double sentimentScore) {
        this.sentimentScore = sentimentScore;
    }

    public Double getConfidenceScore() {
        return confidenceScore;
    }

    public void setConfidenceScore(Double confidenceScore) {
        this.confidenceScore = confidenceScore;
    }

    public String getAnalysis() {
        return analysis;
    }

    public void setAnalysis(String analysis) {
        this.analysis = analysis;
    }

    // Request class for submitting responses
    public static class SubmitResponseRequest {
        @SerializedName("interview_id")
        private String interviewId;

        @SerializedName("question_id")
        private String questionId;

        @SerializedName("response_text")
        private String responseText;

        @SerializedName("duration")
        private int duration;

        public SubmitResponseRequest(String interviewId, String questionId, String responseText, int duration) {
            this.interviewId = interviewId;
            this.questionId = questionId;
            this.responseText = responseText;
            this.duration = duration;
        }
    }

    // Response wrapper class
    public static class SubmitResponseResult {
        @SerializedName("response")
        private Response response;

        @SerializedName("analysis")
        private AnalysisResult analysis;

        public Response getResponse() {
            return response;
        }

        public AnalysisResult getAnalysis() {
            return analysis;
        }
    }

    public static class AnalysisResult {
        @SerializedName("sentiment_score")
        private double sentimentScore;

        @SerializedName("confidence_score")
        private double confidenceScore;

        @SerializedName("feedback")
        private String feedback;

        public double getSentimentScore() {
            return sentimentScore;
        }

        public double getConfidenceScore() {
            return confidenceScore;
        }

        public String getFeedback() {
            return feedback;
        }
    }
}
