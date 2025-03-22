ARG PYTHON=3.9
FROM python:${PYTHON}-slim AS builder

ARG API_VERSION=v1

WORKDIR /builder

# Install pip-tools to compile requirements
RUN pip install --no-cache-dir pip-tools

# Copy all requirements.in files
COPY pipeline/extract/requirements.in ./extractor-requirements.in
COPY common/api/${API_VERSION}/requirements.in ./api-${API_VERSION}-requirements.in
COPY common/utils/requirements.in ./utils-requirements.in

# Combine all requirements.in into one master.in
RUN cat extractor-requirements.in api-${API_VERSION}-requirements.in utils-requirements.in > combined.in

# Compile into final requirements.txt
RUN pip-compile combined.in --output-file=requirements.txt

FROM python:${PYTHON}-slim

ARG API_VERSION=v1

WORKDIR /pythology

RUN mkdir /conf && chown 1000:1000 /conf

COPY --from=builder /builder/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --chown=1000:1000 common/api/${API_VERSION} common/api/${API_VERSION}
COPY --chown=1000:1000 common/utils common/utils
COPY --chown=1000:1000 pipeline/extract/reddit/ reddit/
COPY --chown=1000:1000 pipeline/extract/extractor.py .
COPY --chown=1000:1000 --chmod=755 pipeline/extract/main.py .

USER 1000

CMD ["python3", "main.py"]