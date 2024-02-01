from typing import Any, Callable

from lib.archive import FArchiveReader, FArchiveWriter
from lib.rawdata import (
    base_camp,
    base_camp_module,
    build_process,
    character,
    character_container,
    connector,
    debug,
    dynamic_item,
    foliage_model,
    foliage_model_instance,
    group,
    item_container,
    item_container_slots,
    map_model,
    work,
    work_collection,
    worker_director,
)

PALWORLD_TYPE_HINTS: dict[str, str] = {
    ".worldSaveData.CharacterContainerSaveData.Key": "StructProperty",
    ".worldSaveData.CharacterSaveParameterMap.Key": "StructProperty",
    ".worldSaveData.CharacterSaveParameterMap.Value": "StructProperty",
    ".worldSaveData.FoliageGridSaveDataMap.Key": "StructProperty",
    ".worldSaveData.FoliageGridSaveDataMap.Value.ModelMap.Value": "StructProperty",
    ".worldSaveData.FoliageGridSaveDataMap.Value.ModelMap.Value.InstanceDataMap.Key": "StructProperty",
    ".worldSaveData.FoliageGridSaveDataMap.Value.ModelMap.Value.InstanceDataMap.Value": "StructProperty",
    ".worldSaveData.FoliageGridSaveDataMap.Value": "StructProperty",
    ".worldSaveData.ItemContainerSaveData.Key": "StructProperty",
    ".worldSaveData.MapObjectSaveData.MapObjectSaveData.ConcreteModel.ModuleMap.Value": "StructProperty",
    ".worldSaveData.MapObjectSaveData.MapObjectSaveData.Model.EffectMap.Value": "StructProperty",
    ".worldSaveData.MapObjectSpawnerInStageSaveData.Key": "StructProperty",
    ".worldSaveData.MapObjectSpawnerInStageSaveData.Value": "StructProperty",
    ".worldSaveData.MapObjectSpawnerInStageSaveData.Value.SpawnerDataMapByLevelObjectInstanceId.Key": "Guid",
    ".worldSaveData.MapObjectSpawnerInStageSaveData.Value.SpawnerDataMapByLevelObjectInstanceId.Value": "StructProperty",
    ".worldSaveData.MapObjectSpawnerInStageSaveData.Value.SpawnerDataMapByLevelObjectInstanceId.Value.ItemMap.Value": "StructProperty",
    ".worldSaveData.WorkSaveData.WorkSaveData.WorkAssignMap.Value": "StructProperty",
    ".worldSaveData.BaseCampSaveData.Key": "Guid",
    ".worldSaveData.BaseCampSaveData.Value": "StructProperty",
    ".worldSaveData.BaseCampSaveData.Value.ModuleMap.Value": "StructProperty",
    ".worldSaveData.ItemContainerSaveData.Value": "StructProperty",
    ".worldSaveData.CharacterContainerSaveData.Value": "StructProperty",
    ".worldSaveData.GroupSaveDataMap.Key": "Guid",
    ".worldSaveData.GroupSaveDataMap.Value": "StructProperty",
    ".worldSaveData.EnemyCampSaveData.EnemyCampStatusMap.Value": "StructProperty",
    ".worldSaveData.DungeonSaveData.DungeonSaveData.MapObjectSaveData.MapObjectSaveData.Model.EffectMap.Value": "StructProperty",
    ".worldSaveData.DungeonSaveData.DungeonSaveData.MapObjectSaveData.MapObjectSaveData.ConcreteModel.ModuleMap.Value": "StructProperty",
}

PALWORLD_CUSTOM_PROPERTIES: dict[
    str,
    tuple[
        Callable[[FArchiveReader, str, int, str], dict[str, Any]],
        Callable[[FArchiveWriter, str, dict[str, Any]], int],
    ],
] = {
    ".worldSaveData.GroupSaveDataMap": (group.decode, group.encode),
    ".worldSaveData.CharacterSaveParameterMap.Value.RawData": (
        character.decode,
        character.encode,
    ),
    ".worldSaveData.MapObjectSaveData.MapObjectSaveData.Model.BuildProcess.RawData": (
        build_process.decode,
        build_process.encode,
    ),
    ".worldSaveData.MapObjectSaveData.MapObjectSaveData.Model.Connector.RawData": (
        connector.decode,
        connector.encode,
    ),
    ".worldSaveData.MapObjectSaveData.MapObjectSaveData.Model.RawData": (
        map_model.decode,
        map_model.encode,
    ),
    ".worldSaveData.ItemContainerSaveData.Value.RawData": (
        item_container.decode,
        item_container.encode,
    ),
    ".worldSaveData.ItemContainerSaveData.Value.Slots.Slots.RawData": (
        item_container_slots.decode,
        item_container_slots.encode,
    ),
    # This isn't actually serialised into at all?
    # ".worldSaveData.CharacterContainerSaveData.Value.RawData": (debug.decode, debug.encode),
    # This duplicates the data already serialised into the Slots UObject?
    ".worldSaveData.CharacterContainerSaveData.Value.Slots.Slots.RawData": (
        character_container.decode,
        character_container.encode,
    ),
    # DynamicItemSaveData is problematic because serialisation is dependent on type, which is not immediately obvious
    ".worldSaveData.DynamicItemSaveData.DynamicItemSaveData.RawData": (
        dynamic_item.decode,
        dynamic_item.encode,
    ),
    ".worldSaveData.FoliageGridSaveDataMap.Value.ModelMap.Value.RawData": (
        foliage_model.decode,
        foliage_model.encode,
    ),
    ".worldSaveData.FoliageGridSaveDataMap.Value.ModelMap.Value.InstanceDataMap.Value.RawData": (
        foliage_model_instance.decode,
        foliage_model_instance.encode,
    ),
    ".worldSaveData.BaseCampSaveData.Value.RawData": (
        base_camp.decode,
        base_camp.encode,
    ),
    ".worldSaveData.BaseCampSaveData.Value.WorkerDirector.RawData": (
        worker_director.decode,
        worker_director.encode,
    ),
    ".worldSaveData.BaseCampSaveData.Value.WorkCollection.RawData": (
        work_collection.decode,
        work_collection.encode,
    ),
    # ".worldSaveData.BaseCampSaveData.Value.ModuleMap": (base_camp_module.decode, base_camp_module.encode),
    ".worldSaveData.WorkSaveData": (work.decode, work.encode),
    # ".worldSaveData.WorkSaveData.WorkSaveData.RawData": (debug.decode, debug.encode),
    # ".worldSaveData.WorkSaveData.WorkSaveData.WorkAssignMap.Value.RawData": (debug.decode, debug.encode),
    # ConcreteModel is problematic because serialisation is dependent on type, which is not immediately obvious
    # ".worldSaveData.MapObjectSaveData.MapObjectSaveData.ConcreteModel": (
    #     decode_map_concrete_model,
    #     encode_map_concrete_model,
    # ),
    # ".worldSaveData.MapObjectSaveData.MapObjectSaveData.ConcreteModel.RawData": (),
    # ".worldSaveData.MapObjectSaveData.MapObjectSaveData.ConcreteModel.ModuleMap.Value.RawData": (),
}
