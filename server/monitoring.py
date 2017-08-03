# Monitoring instrumentation module for Flask web servers. Thread-safe.
# It supports simple counters.

from sys import maxsize
from threading import RLock

class Metric(object):
  """Metric represents a single metric."""

  def __init__(self, name):
    self.lock = RLock()
    self.name = name
    self._value = 0

  def increment(self):
    with self.lock:
      if self._value >= maxsize:
        # don't store counter values that are greater than it would be practical
        # reset cnt, same happens on restart (after crash)
        # side effect: some intervals will be inaccurate
        self.reset()
      self._value += 1

  def reset(self):
    with self.lock:
      self._value = 0

  def get(self):
    with self.lock:
      return self._value

def make_lambda(metric):
  """Simulate a closure (helper due to Python 2's lack of closure support)."""
  return lambda: "%d" % metric.get()

class MetricStore(object):
  """MetricStore represents a collection of metrics."""

  def __init__(self, *names):
    metrics = dict()
    for name in names:
      metrics[name] = Metric(name)
    self._metrics = metrics

  def __call__(self, name):
    return self._metrics[name]

  def export_to(self, app):
    """Export metric endpoints to Flask app"""
    for metric in self._metrics.values():
      app.add_url_rule('/metrics/%s' % metric.name, '/metrics/%s' % metric.name, make_lambda(metric))

def create_metrics(*names):
  """Create a MetricStore."""
  return MetricStore(*names)
