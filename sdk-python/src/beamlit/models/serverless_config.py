from typing import Any, Type, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="ServerlessConfig")


@_attrs_define
class ServerlessConfig:
    """Configuration for the serverless model deployment

    Attributes:
        last_pod_retention_period (Union[Unset, str]): The last pod retention period
        max_num_replicas (Union[Unset, int]): The maximum number of replicas
        metric (Union[Unset, str]): Metric to scale on, can be "cpu" or "memory" or "rps" or "concurrency"
        min_num_replicas (Union[Unset, int]): The minimum number of replicas
        scale_down_delay (Union[Unset, str]): The scale down delay
        scale_up_minimum (Union[Unset, int]): The scale up minimum
        stable_window (Union[Unset, str]): The stable window
        target (Union[Unset, str]): Target value for the metric
    """

    last_pod_retention_period: Union[Unset, str] = UNSET
    max_num_replicas: Union[Unset, int] = UNSET
    metric: Union[Unset, str] = UNSET
    min_num_replicas: Union[Unset, int] = UNSET
    scale_down_delay: Union[Unset, str] = UNSET
    scale_up_minimum: Union[Unset, int] = UNSET
    stable_window: Union[Unset, str] = UNSET
    target: Union[Unset, str] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        last_pod_retention_period = self.last_pod_retention_period

        max_num_replicas = self.max_num_replicas

        metric = self.metric

        min_num_replicas = self.min_num_replicas

        scale_down_delay = self.scale_down_delay

        scale_up_minimum = self.scale_up_minimum

        stable_window = self.stable_window

        target = self.target

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update({})
        if last_pod_retention_period is not UNSET:
            field_dict["last_pod_retention_period"] = last_pod_retention_period
        if max_num_replicas is not UNSET:
            field_dict["max_num_replicas"] = max_num_replicas
        if metric is not UNSET:
            field_dict["metric"] = metric
        if min_num_replicas is not UNSET:
            field_dict["min_num_replicas"] = min_num_replicas
        if scale_down_delay is not UNSET:
            field_dict["scale_down_delay"] = scale_down_delay
        if scale_up_minimum is not UNSET:
            field_dict["scale_up_minimum"] = scale_up_minimum
        if stable_window is not UNSET:
            field_dict["stable_window"] = stable_window
        if target is not UNSET:
            field_dict["target"] = target

        return field_dict

    @classmethod
    def from_dict(cls: Type[T], src_dict: dict[str, Any]) -> T:
        d = src_dict.copy()
        last_pod_retention_period = d.pop("last_pod_retention_period", UNSET)

        max_num_replicas = d.pop("max_num_replicas", UNSET)

        metric = d.pop("metric", UNSET)

        min_num_replicas = d.pop("min_num_replicas", UNSET)

        scale_down_delay = d.pop("scale_down_delay", UNSET)

        scale_up_minimum = d.pop("scale_up_minimum", UNSET)

        stable_window = d.pop("stable_window", UNSET)

        target = d.pop("target", UNSET)

        serverless_config = cls(
            last_pod_retention_period=last_pod_retention_period,
            max_num_replicas=max_num_replicas,
            metric=metric,
            min_num_replicas=min_num_replicas,
            scale_down_delay=scale_down_delay,
            scale_up_minimum=scale_up_minimum,
            stable_window=stable_window,
            target=target,
        )

        serverless_config.additional_properties = d
        return serverless_config

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
