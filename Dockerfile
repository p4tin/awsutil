FROM unbounce/yopa

EXPOSE 47195
EXPOSE 47196
EXPOSE 47197

ADD ./config /config
CMD java -Xms64m -Xmx256m -jar uberjar.jar -c /config/yopa.yaml

