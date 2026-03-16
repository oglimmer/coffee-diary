package com.coffeediary.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Table(name = "diary_entries")
public class DiaryEntry {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Column(name = "date_time", nullable = false)
    private LocalDateTime dateTime;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "sieve_id")
    private Sieve sieve;

    @Column(nullable = false)
    @Builder.Default
    private Integer temperature = 93;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "coffee_id")
    private Coffee coffee;

    @Column(name = "grind_size")
    private Double grindSize;

    @Column(name = "input_weight")
    private Double inputWeight;

    @Column(name = "output_weight")
    private Double outputWeight;

    @Column(name = "time_seconds")
    private Integer timeSeconds;

    private Integer rating;

    @Column(columnDefinition = "TEXT")
    private String notes;
}
