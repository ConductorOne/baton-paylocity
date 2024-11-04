FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-paylocity"]
COPY baton-paylocity /