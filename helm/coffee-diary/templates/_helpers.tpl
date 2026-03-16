{{/*
Expand the name of the chart.
*/}}
{{- define "coffee-diary.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "coffee-diary.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "coffee-diary.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "coffee-diary.labels" -}}
helm.sh/chart: {{ include "coffee-diary.chart" . }}
{{ include "coffee-diary.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "coffee-diary.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coffee-diary.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Backend specific labels
*/}}
{{- define "coffee-diary.backend.labels" -}}
helm.sh/chart: {{ include "coffee-diary.chart" . }}
{{ include "coffee-diary.backend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Backend selector labels
*/}}
{{- define "coffee-diary.backend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coffee-diary.name" . }}-backend
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: backend
{{- end }}

{{/*
Frontend specific labels
*/}}
{{- define "coffee-diary.frontend.labels" -}}
helm.sh/chart: {{ include "coffee-diary.chart" . }}
{{ include "coffee-diary.frontend.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Frontend selector labels
*/}}
{{- define "coffee-diary.frontend.selectorLabels" -}}
app.kubernetes.io/name: {{ include "coffee-diary.name" . }}-frontend
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: frontend
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "coffee-diary.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "coffee-diary.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Database host
*/}}
{{- define "coffee-diary.database.host" -}}
{{- .Values.database.host }}
{{- end }}

{{/*
Database port
*/}}
{{- define "coffee-diary.database.port" -}}
{{- .Values.database.port }}
{{- end }}

{{/*
Database name
*/}}
{{- define "coffee-diary.database.name" -}}
{{- .Values.database.name }}
{{- end }}

{{/*
Backend fullname
*/}}
{{- define "coffee-diary.backend.fullname" -}}
{{- printf "%s-backend" (include "coffee-diary.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Frontend fullname
*/}}
{{- define "coffee-diary.frontend.fullname" -}}
{{- printf "%s-frontend" (include "coffee-diary.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}
