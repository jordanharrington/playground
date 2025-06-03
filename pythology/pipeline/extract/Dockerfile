ARG PYTHON=3.9
FROM python:${PYTHON}-slim AS builder

WORKDIR /builder

RUN pip install --no-cache-dir pip-tools

COPY pipeline/extract/requirements.in ./requirements/extractor.in
COPY common/api/requirements.in ./requirements/api.in
COPY common/utils/requirements.in ./requirements/utils.in

RUN pip-compile requirements/*.in --output-file=requirements.txt --strip-extras

FROM python:${PYTHON}-slim

WORKDIR /pythology

RUN mkdir /conf && chown 1000:1000 /conf

COPY --from=builder /builder/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY --chown=1000:1000 --chmod=755 common/ common/
COPY --chown=1000:1000 --chmod=755 pipeline/extract/ pipeline/extract

ENV PYTHONPATH=/pythology

USER 1000

ENTRYPOINT ["python", "pipeline/extract/main.py"]