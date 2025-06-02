ARG PYTHON=3.9
FROM python:${PYTHON}-slim AS builder

ARG API_VERSION=v1

WORKDIR /builder

RUN pip install --no-cache-dir pip-tools

COPY pipeline/extract/requirements.in ./requirements/extractor.in
COPY common/api/${API_VERSION}/requirements.in ./requirements/api.in
COPY common/utils/requirements.in ./requirements/utils.in

RUN pip-compile requirements/*.in --output-file=requirements.txt --strip-extras

FROM python:${PYTHON}-slim

ARG API_VERSION=v1

WORKDIR /pythology

RUN mkdir /conf && chown 1000:1000 /conf

COPY --from=builder /builder/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --chown=1000:1000 common/api/${API_VERSION} common/api/${API_VERSION}
COPY --chown=1000:1000 common/utils common/utils
COPY --chown=1000:1000 pipeline/extract/reddit reddit/
COPY --chown=1000:1000 pipeline/extract/extractor.py .
COPY --chown=1000:1000 --chmod=755 pipeline/extract/main.py .

USER 1000

ENTRYPOINT ["python", "main.py"]