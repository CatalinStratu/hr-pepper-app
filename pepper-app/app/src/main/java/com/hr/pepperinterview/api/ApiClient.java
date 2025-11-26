package com.hr.pepperinterview.api;

import com.hr.pepperinterview.BuildConfig;

import java.util.concurrent.TimeUnit;

import okhttp3.OkHttpClient;
import okhttp3.logging.HttpLoggingInterceptor;
import retrofit2.Retrofit;
import retrofit2.converter.gson.GsonConverterFactory;

public class ApiClient {
    private static ApiClient instance;
    private final ApiService apiService;
    private String serverUrl;

    private ApiClient(String serverUrl) {
        this.serverUrl = serverUrl;

        HttpLoggingInterceptor logging = new HttpLoggingInterceptor();
        logging.setLevel(HttpLoggingInterceptor.Level.BODY);

        OkHttpClient client = new OkHttpClient.Builder()
                .addInterceptor(logging)
                .connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .writeTimeout(30, TimeUnit.SECONDS)
                .build();

        Retrofit retrofit = new Retrofit.Builder()
                .baseUrl(serverUrl)
                .client(client)
                .addConverterFactory(GsonConverterFactory.create())
                .build();

        apiService = retrofit.create(ApiService.class);
    }

    public static synchronized ApiClient getInstance() {
        if (instance == null) {
            instance = new ApiClient(BuildConfig.SERVER_URL);
        }
        return instance;
    }

    public static synchronized void initialize(String serverUrl) {
        instance = new ApiClient(serverUrl);
    }

    public ApiService getApiService() {
        return apiService;
    }

    public String getServerUrl() {
        return serverUrl;
    }
}
