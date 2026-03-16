package com.coffeediary.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class DiaryEntryResponse {

    private Long id;
    private Long userId;
    private LocalDateTime dateTime;
    private Long sieveId;
    private String sieveName;
    private Integer temperature;
    private Long coffeeId;
    private String coffeeName;
    private Double grindSize;
    private Double inputWeight;
    private Double outputWeight;
    private Integer timeSeconds;
    private Integer rating;
    private String notes;
}
